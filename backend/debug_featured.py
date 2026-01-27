import urllib.request
import json
import sys

try:
    url = 'http://localhost:8080/api/featured'
    print(f"Fetching {url}...")
    
    with urllib.request.urlopen(url) as response:
        content = response.read().decode()
        data = json.loads(content)
    
    print(f"Status: {response.status}")
    print(f"Items: {len(data)}")
    
    if len(data) > 0:
        item = data[0]
        print("\nSample Item Keys:", item.keys())
        
        print(f"ID: {item.get('id')}")
        print(f"ModelID: {item.get('model_id')}")
        
        model = item.get('model')
        if model:
            print(f"Model found: {model.get('name')}")
            if model.get('brand'):
                print(f"Brand found: {model.get('brand', {}).get('name')}")
            else:
                print("Brand MISSING in model!")
        else:
            print("Model MISSING in item!")
            
    # Check for nulls in critical fields
    for i, item in enumerate(data):
        if not item.get('model'):
             print(f"WARNING: Item {i} (ID: {item.get('id')}) has NO MODEL!")
        elif not item.get('model', {}).get('brand'):
             print(f"WARNING: Item {i} (ID: {item.get('id')}) has NO BRAND via Model!")
             
except Exception as e:
    print(f"ERROR: {e}")
    # Print raw content if json parse fails
    # print(content[:200])
