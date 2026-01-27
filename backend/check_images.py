import sqlite3

conn = sqlite3.connect('backend/vehicles.db')
cursor = conn.cursor()

# Check images
cursor.execute("SELECT COUNT(*) FROM trims WHERE image_url IS NOT NULL AND image_url != ''")
with_images = cursor.fetchone()[0]

cursor.execute("SELECT COUNT(*) FROM trims")
total = cursor.fetchone()[0]

print(f"Trims with images: {with_images}/{total}")
print()

if with_images == 0:
    print("⚠ NO TRIMS HAVE IMAGES!")
    print("This is why /api/featured returns null - it filters by image_url")
    print()
    print("Sample trims (first 5):")
    cursor.execute("SELECT id, name, year, image_url FROM trims LIMIT 5")
    for row in cursor.fetchall():
        img_status = "HAS IMAGE" if row[3] else "NO IMAGE"
        print(f"  {row[0]}: {row[1]} ({row[2]}) - {img_status}")

conn.close()
