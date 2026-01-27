import sqlite3
import sys

def add_column_if_not_exists(cursor, table, column, definition):
    """Add column to table if it doesn't already exist"""
    cursor.execute(f"PRAGMA table_info({table})")
    columns = [row[1] for row in cursor.fetchall()]
    
    if column not in columns:
        print(f"  Adding column {column} to {table}...")
        cursor.execute(f"ALTER TABLE {table} ADD COLUMN {column} {definition}")
        return True
    else:
        print(f"  Column {column} already exists in {table}")
        return False

def run_safe_migration():
    """Run safe 4-level migration without recreating tables"""
    
    conn = sqlite3.connect('backend/vehicles.db')
    cursor = conn.cursor()
    
    print("=" * 70)
    print("🚀 SAFE 4-LEVEL MIGRATION")
    print("=" * 70)
    print()
    
    try:
        # VERIFY STARTING STATE
        print("📊 Pre-migration verification:")
        cursor.execute("SELECT COUNT(*) FROM trims")
        initial_trims = cursor.fetchone()[0]
        print(f"  Trims before migration: {initial_trims}")
        
        cursor.execute("SELECT COUNT(*) FROM trims WHERE image_url IS NOT NULL")
        trims_with_images = cursor.fetchone()[0]
        print(f"  Trims with images: {trims_with_images}")
        print()
        
        if initial_trims == 0:
            print("❌  ERROR: No trims found in database!")
            return False
        
        # STEP 1: Create generations table
        print("Step 1: Creating generations table...")
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS generations (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                model_id INTEGER NOT NULL,
                code TEXT NOT NULL,
                name TEXT,
                start_year INTEGER NOT NULL,
                end_year INTEGER,
                image_url TEXT, description TEXT,
                is_current BOOLEAN DEFAULT 0,
                platform TEXT,
                created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (model_id) REFERENCES models(id) ON DELETE CASCADE
            )
        """)
        print("  ✓ Generations table created")
        
        # STEP 2: Populate generations
        print("\nStep 2: Populating generations from trims...")
        cursor.execute("""
            INSERT INTO generations (model_id, code, name, start_year, end_year, is_current)
            SELECT DISTINCT
                model_id,
                code,
                code || ' (' || CAST(min_year AS TEXT) || 
                    CASE WHEN max_year != min_year THEN '-' || CAST(max_year AS TEXT) ELSE '' END || ')' as name,
                min_year as start_year,
                CASE WHEN max_year >= 2024 THEN NULL ELSE max_year END as end_year,
                CASE WHEN max_year >= 2024 THEN 1 ELSE 0 END as is_current
            FROM (
                SELECT
                    t.model_id,
                    COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT)) as code,
                    MIN(t.year) as min_year,
                    MAX(t.year) as max_year
                FROM trims t
                WHERE t.model_id IS NOT NULL
                GROUP BY t.model_id, COALESCE(NULLIF(t.generation, ''), 'Y' || CAST(t.year AS TEXT))
            )
        """)
        
        cursor.execute("SELECT COUNT(*) FROM generations")
        gen_count = cursor.fetchone()[0]
        print(f"  ✓ Created {gen_count} generations")
        
        # STEP 3: Add generation_id column
        print("\nStep 3: Adding generation_id column to trims...")
        added = add_column_if_not_exists(cursor, 'trims', 'generation_id', 'INTEGER')
        
        # STEP 4: Update generation_id for all trims
        print("\nStep 4: Linking trims to generations...")
        cursor.execute("""
            UPDATE trims
            SET generation_id = (
                SELECT g.id
                FROM generations g
                WHERE g.model_id = trims.model_id
                  AND g.code = COALESCE(NULLIF(trims.generation, ''), 'Y' || CAST(trims.year AS TEXT))
                  AND trims.year >= g.start_year
                  AND (g.end_year IS NULL OR trims.year <= g.end_year)
                LIMIT 1
            )
        """)
        print("  ✓ Trims linked to generations")
        
        # STEP 5: Create indexes
        print("\nStep 5: Creating indexes...")
        cursor.execute("CREATE INDEX IF NOT EXISTS idx_generations_model ON generations(model_id)")
        cursor.execute("CREATE INDEX IF NOT EXISTS idx_generations_code ON generations(code)")
        cursor.execute("CREATE INDEX IF NOT EXISTS idx_trims_generation ON trims(generation_id)")
        print("  ✓ Indexes created")
        
        conn.commit()
        
        # VERIFICATION
        print("\n" + "=" * 70)
        print("✅ MIGRATION VERIFICATION")
        print("=" * 70)
        
        cursor.execute("SELECT COUNT(*) FROM trims")
        final_trims = cursor.fetchone()[0]
        
        cursor.execute("SELECT COUNT(*) FROM trims WHERE generation_id IS NOT NULL")
        linked_trims = cursor.fetchone()[0]
        
        cursor.execute("SELECT COUNT(*) FROM trims WHERE generation_id IS NULL")
        orphaned = cursor.fetchone()[0]
        
        cursor.execute("SELECT COUNT(*) FROM generations")
        total_gens = cursor.fetchone()[0]
        
        print(f"  Total trims: {final_trims} (started with {initial_trims})")
        print(f"  Trims linked to generations: {linked_trims}")
        print(f"  Orphaned trims: {orphaned}")
        print(f"  Total generations: {total_gens}")
        print()
        
        if final_trims != initial_trims:
            print(f"❌ ERROR: Data loss detected! {initial_trims - final_trims} trims lost!")
            conn.rollback()
            return False
        
        if orphaned > 0:
            print(f"⚠ WARNING: {orphaned} trims not linked to generations")
        
        print("=" * 70)
        print("✅ MIGRATION SUCCESSFUL!")
        print("=" * 70)
        return True
        
    except Exception as e:
        print(f"\n❌ Migration failed: {e}")
        import traceback
        traceback.print_exc()
        conn.rollback()
        return False
    finally:
        conn.close()

if __name__ == "__main__":
    success = run_safe_migration()
    sys.exit(0 if success else 1)
