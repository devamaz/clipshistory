package store

import (
	"crypto/sha256"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Clip struct {
	ID           int64
	Content      string
	ContentHash  string
	Preview      string
	CharCount    int
	CreatedAt    int64
	LastCopiedAt int64
	CopyCount    int
	IsPinned     bool
	IsDeleted    bool
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite", "./aclips")
	if err != nil {
		log.Fatal(err)
	}

	createTableStmt := `
		CREATE TABLE IF NOT EXISTS aclips (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT NOT NULL,
				content_hash TEXT NOT NULL UNIQUE,
				preview TEXT NOT NULL,
				char_count INTEGER NOT NULL,
				created_at INTEGER NOT NULL,
				last_copied_at INTEGER NOT NULL,
				copy_count INTEGER DEFAULT 1,
				is_pinned BOOLEAN DEFAULT 0,
				is_deleted BOOLEAN DEFAULT 0
		);

		-- Indexes for performance
		CREATE INDEX IF NOT EXISTS idx_last_copied ON aclips(last_copied_at DESC);
		CREATE INDEX IF NOT EXISTS idx_pinned_date ON aclips(is_pinned DESC, last_copied_at DESC);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_content_hash ON aclips(content_hash);
		CREATE INDEX IF NOT EXISTS idx_deleted ON aclips(is_deleted);

		-- Optional: FTS5 for blazing fast search
		CREATE VIRTUAL TABLE IF NOT EXISTS clipboard_search USING fts5(
				content,
				content=aclips,
				content_rowid=id
		);

		-- Triggers to keep FTS in sync (if using FTS5)
		CREATE TRIGGER IF NOT EXISTS history_ai AFTER INSERT ON aclips BEGIN
				INSERT INTO clipboard_search(rowid, content) VALUES (new.id, new.content);
		END;

		CREATE TRIGGER IF NOT EXISTS history_ad AFTER DELETE ON aclips BEGIN
				DELETE FROM clipboard_search WHERE rowid = old.id;
		END;

		CREATE TRIGGER IF NOT EXISTS history_au AFTER UPDATE ON aclips BEGIN
				UPDATE clipboard_search SET content = new.content WHERE rowid = new.id;
		END;
	`

	if _, err := s.conn.Exec(createTableStmt); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Store) GetClips() ([]Clip, error) {
	rows, err := s.conn.Query("SELECT * FROM aclips WHERE is_deleted = 0")
	if err != nil {
		return nil, err
	}

	clips := []Clip{}
	defer rows.Close()
	for rows.Next() {
		clip := Clip{}
		rows.Scan(&clip.ID, &clip.Content, &clip.ContentHash, &clip.Preview, &clip.CharCount, &clip.CreatedAt, &clip.LastCopiedAt, &clip.CopyCount, &clip.IsPinned, &clip.IsDeleted)
		clips = append(clips, clip)
	}

	return clips, nil
}

func (s *Store) GetPinnedClips() ([]Clip, error) {
	rows, err := s.conn.Query("SELECT * FROM aclips WHERE is_pinned = 1 AND is_deleted = 0")
	if err != nil {
		return nil, err
	}

	clips := []Clip{}
	defer rows.Close()
	for rows.Next() {
		clip := Clip{}
		rows.Scan(&clip.ID, &clip.Content, &clip.ContentHash, &clip.Preview, &clip.CharCount, &clip.CreatedAt, &clip.LastCopiedAt, &clip.CopyCount, &clip.IsPinned, &clip.IsDeleted)
		clips = append(clips, clip)
	}

	return clips, nil
}

func (a *Store) GetClip(id int) (Clip, error) {
	var clip Clip
	err := a.conn.QueryRow("SELECT * FROM aclips WHERE id = ?", id).Scan(&clip.ID, &clip.Preview, &clip.Content, &clip.CharCount, &clip.CreatedAt, &clip.LastCopiedAt, &clip.CopyCount, &clip.IsPinned, &clip.IsDeleted)
	if err != nil {
		return clip, err
	}

	return clip, nil
}

func (s *Store) Save(clip Clip) error {

	contentToHash := clip.Content
	h := sha256.New()              // Create a new SHA-256 hash object
	h.Write([]byte(contentToHash)) // Write the input string as bytes to the hash object
	ContentHash := h.Sum(nil)

	charCount := len(clip.Content)

	preview := clip.Content[:min(len(clip.Content), 15)]

	upsertQuery := `INSERT INTO aclips ( content, content_hash, preview, char_count, created_at, last_copied_at, copy_count, is_pinned, is_deleted)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?,?)
	ON CONFLICT(content_hash) DO UPDATE SET
	last_copied_at=excluded.last_copied_at, copy_count=copy_count+1;`

	if _, err := s.conn.Exec(upsertQuery, clip.Content, ContentHash, preview, charCount, clip.CreatedAt, clip.LastCopiedAt, clip.CopyCount, clip.IsPinned, clip.IsDeleted); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteNote(id int) error {
	deleteQuery := `DELETE FROM notes WHERE id = ?`

	if _, err := s.conn.Exec(deleteQuery, id); err != nil {
		return err
	}

	return nil
}

// func (s *Store) SearchClip(searchTerm string, useFTS bool) ([]Clip, error) {
// 	searchTerm = strings.TrimSpace(searchTerm)

// 	if useFTS {
// 		return searchWithFTS(db, searchTerm)
// 	}

// 	return searchWithLike(db, searchTerm)
// }

// // searchWithFTS performs full-text search using FTS5
// func searchWithFTS(db *sql.DB, searchTerm string) ([]ClipboardEntry, error) {
// 	// Convert spaces to AND for multi-word search
// 	// "python function" becomes "python AND function"
// 	words := strings.Fields(searchTerm)
// 	ftsQuery := strings.Join(words, " AND ")

// 	query := `
//         SELECT h.id, h.preview, h.char_count, h.last_copied_at, h.copy_count, h.is_pinned
//         FROM clipboard_search s
//         JOIN clipboard_history h ON h.id = s.rowid
//         WHERE clipboard_search MATCH ?
//           AND h.is_deleted = 0
//         ORDER BY h.is_pinned DESC, h.copy_count DESC, h.last_copied_at DESC
//         LIMIT 50
//     `

// 	return executeQuery(db, query, ftsQuery)
// }

// // searchWithLike performs basic substring search
// func searchWithLike(db *sql.DB, searchTerm string) ([]ClipboardEntry, error) {
// 	query := `
//         SELECT id, preview, char_count, last_copied_at, copy_count, is_pinned
//         FROM clipboard_history
//         WHERE is_deleted = 0
//           AND content LIKE ? COLLATE NOCASE
//         ORDER BY is_pinned DESC, copy_count DESC, last_copied_at DESC
//         LIMIT 50
//     `

// 	likePattern := fmt.Sprintf("%%%s%%", searchTerm)
// 	return executeQuery(db, query, likePattern)
// }

// // executeQuery is a helper to run queries and scan results
// func executeQuery(db *sql.DB, query string, args ...interface{}) ([]ClipboardEntry, error) {
// 	rows, err := db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("query failed: %w", err)
// 	}
// 	defer rows.Close()

// 	var entries []Clip

// 	for rows.Next() {
// 		var entry ClipboardEntry
// 		var isPinned int // SQLite stores booleans as integers

// 		err := rows.Scan(
// 			&entry.ID,
// 			&entry.Preview,
// 			&entry.CharCount,
// 			&entry.LastCopiedAt,
// 			&entry.CopyCount,
// 			&isPinned,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan failed: %w", err)
// 		}

// 		entry.IsPinned = isPinned == 1
// 		entries = append(entries, entry)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows iteration failed: %w", err)
// 	}

// 	return entries, nil
// }

// // SearchWithPrefixMatch performs prefix matching (search-as-you-type with FTS5)
// func SearchWithPrefixMatch(db *sql.DB, searchTerm string) ([]ClipboardEntry, error) {
// 	searchTerm = strings.TrimSpace(searchTerm)

// 	if searchTerm == "" {
// 		return getAllEntries(db)
// 	}

// 	// Add * suffix for prefix matching
// 	ftsQuery := searchTerm + "*"

// 	query := `
//         SELECT h.id, h.preview, h.char_count, h.last_copied_at, h.copy_count, h.is_pinned
//         FROM clipboard_search s
//         JOIN clipboard_history h ON h.id = s.rowid
//         WHERE clipboard_search MATCH ?
//           AND h.is_deleted = 0
//         ORDER BY h.is_pinned DESC, h.last_copied_at DESC
//         LIMIT 50
//     `

// 	return executeQuery(db, query, ftsQuery)
// }

// // FormatTimestamp converts Unix timestamp to human-readable format
// func FormatTimestamp(timestamp int64) string {
// 	t := time.Unix(timestamp, 0)

// 	now := time.Now()
// 	diff := now.Sub(t)

// 	switch {
// 	case diff < time.Minute:
// 		return "just now"
// 	case diff < time.Hour:
// 		mins := int(diff.Minutes())
// 		return fmt.Sprintf("%d min ago", mins)
// 	case diff < 24*time.Hour:
// 		hours := int(diff.Hours())
// 		return fmt.Sprintf("%d hours ago", hours)
// 	case diff < 7*24*time.Hour:
// 		days := int(diff.Hours() / 24)
// 		return fmt.Sprintf("%d days ago", days)
// 	default:
// 		return t.Format("Jan 02, 2006")
// 	}
// }

// // Example usage
// func main() {
// 	db, err := sql.Open("sqlite3", "./clipboard.db")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	// Search with FTS5
// 	results, err := SearchClipboard(db, "python function", true)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Found %d entries:\n\n", len(results))

// 	for i, entry := range results {
// 		pinned := ""
// 		if entry.IsPinned {
// 			pinned = "📌 "
// 		}

// 		fmt.Printf("%d. %s%s\n", i+1, pinned, entry.Preview)
// 		fmt.Printf("   Copied %d times | %s | %d chars\n\n",
// 			entry.CopyCount,
// 			FormatTimestamp(entry.LastCopiedAt),
// 			entry.CharCount,
// 		)
// 	}
// }
