import sqlite3
import os
import glob

DB_PATH = 'vehicles.db'
MIGRATIONS_DIR = 'migrations'

def run_migrations():
    """Apply all SQL migrations in order"""
    try:
        conn = sqlite3.connect(DB_PATH, timeout=10.0)
        cursor = conn.cursor()
        
        # Enable foreign keys
        cursor.execute("PRAGMA foreign_keys = ON")
        
        # Get all migration files
        migration_files = sorted(glob.glob(os.path.join(MIGRATIONS_DIR, '*.sql')))
        
        if not migration_files:
            print("No migration files found.")
            return
        
        print("Running migrations...")
        for migration_file in migration_files:
            print(f"Applying: {os.path.basename(migration_file)}")
            
            with open(migration_file, 'r', encoding='utf-8') as f:
                sql = f.read()
            
            cursor.executescript(sql)
        
        conn.commit()
        print("\n✓ All migrations applied successfully")
        
        # Verify tables were created
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
        tables = cursor.fetchall()
        print(f"\nCreated tables: {', '.join([t[0] for t in tables])}")
        
    except sqlite3.Error as e:
        print(f"Database error: {e}")
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    run_migrations()
