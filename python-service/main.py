import os
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain_community.tools.tavily_search import TavilySearchResults
from langchain.agents import create_tool_calling_agent, AgentExecutor
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.messages import SystemMessage

# Load environment variables
load_dotenv()

# Initialize FastAPI
app = FastAPI(title="Car Spec AI Microservice")

# Configuration
GOOGLE_API_KEY = os.getenv("GOOGLE_API_KEY")
TAVILY_API_KEY = os.getenv("TAVILY_API_KEY")

if not GOOGLE_API_KEY:
    print("Warning: GOOGLE_API_KEY not found in environment")
if not TAVILY_API_KEY:
    print("Warning: TAVILY_API_KEY not found in environment")

# Initialize LLM (Gemini 1.5 Flash)
llm = ChatGoogleGenerativeAI(
    model="gemini-1.5-flash",
    temperature=0, # Reduce hallucinations
    google_api_key=GOOGLE_API_KEY
)

# Initialize Tools
tavily_tool = TavilySearchResults(k=3)
tools = [tavily_tool]

# Define System Prompt
SYSTEM_PROMPT = """SEN: Kıdemli bir Otomotiv Mühendisi ve Servis Teknisyeni uzmanısın. Adın 'Car Spec AI'.

KURALLAR:
1. **DOĞRULUĞUNDAN EMİN OLDUĞUN BİLGİLERİ KULLAN (ÇOK ÖNEMLİ):**
   - Cevap vermeden önce aracın Benzinli mi Dizel mi olduğunu analiz et.
   - **NEGATİF KISITLAMA:** Benzinli (TSI, TFSI, PureTech, TCe) motorlar için ASLA "DPF", "Kızdırma Bujisi" (Glow Plug) veya "Enjektör İşemesi" (Dizel terimleri) deme.
   - Benzinli araçlarda "Katalitik Konvertör", "Buji", "Bobin" veya "OPF" (Partikül Filtresi - Benzinli) bulunur. Yanlış terminoloji kullanma.

2. **ZORUNLU TOOL KULLANIMI:**
   - Şu konularda soru gelirse KESİNLİKLE `tavily_search_results_json` aracını kullan:
     - "Kronik sorunlar"
     - "Kullanıcı yorumları"
     - "Gerçek yakıt tüketimi"
     - "Güncel piyasa fiyatları"
   - İçsel bilgine güvenme, internetten teyit et.

3. **FORMAT:**
   - Cevaplarını Türkçe ver.
   - Okunabilir bir Markdown listesi formatında sun.
   - Teknik terimleri (Örn: Timing Chain -> Zincir) parantez içinde açıkla.

4. **TON:**
   - Profesyonel, teknik, net ve yardımsever.

5. **BİLİNMEYEN DURUMLAR:**
   - Emin değilsen veya internette bulamazsan "Bu konuda net veri bulamadım" de, uydurma.
"""

# Create Prompt Template
prompt = ChatPromptTemplate.from_messages([
    ("system", SYSTEM_PROMPT),
    ("human", "{input}"),
    ("placeholder", "{agent_scratchpad}"),
])

# Create Agent
agent = create_tool_calling_agent(llm, tools, prompt)
agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

# Pydantic Model for Input
class QueryRequest(BaseModel):
    message: str

# API Endpoints
@app.post("/ask")
async def ask_question(request: QueryRequest):
    try:
        response = agent_executor.invoke({"input": request.message})
        return {"response": response["output"]}
    except Exception as e:
        print(f"Error processing request: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "ok", "service": "Car Spec AI Python Service"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
