import sqlite3

# Connect to database
conn = sqlite3.connect('vehicles.db')
cursor = conn.cursor()

# Check if column exists
cursor.execute("PRAGMA table_info(trims)")
columns = [row[1] for row in cursor.fetchall()]

if 'image_url' not in columns:
    print("Adding image_url column to trims table...")
    cursor.execute("ALTER TABLE trims ADD COLUMN image_url TEXT")
    conn.commit()
    print("✓ Column added successfully!")
else:
    print("✓ image_url column already exists")

conn.close()
