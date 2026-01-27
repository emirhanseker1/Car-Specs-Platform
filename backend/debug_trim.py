import urllib.request
import json
import sys

try:
    url = 'http://localhost:8080/api/trims/39?include_relations=true'
    with urllib.request.urlopen(url) as response:
        data = json.loads(response.read().decode())
    
    print("=" * 50)
    print(f"TRIM 39 DATA CHECK")
    print("=" * 50)
    
    print(f"ID: {data.get('id')}")
    print(f"GenerationID: {data.get('generation_id')}")
    print(f"ModelID: {data.get('model_id')} (Legacy)")
    
    model = data.get('model')
    if model:
        print(f"Model: FOUND - ID: {model.get('id')}, Name: {model.get('name')}")
        brand = model.get('brand')
        if brand:
            print(f"Brand: FOUND - ID: {brand.get('id')}, Name: {brand.get('name')}")
        else:
            print("Brand: MISSING!")
    else:
        print("Model: MISSING!")
        
    print("\nFull JSON:")
    print(json.dumps(data, indent=2))
    
except Exception as e:
    print(f"Error: {e}")
