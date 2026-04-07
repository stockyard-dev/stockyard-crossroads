package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type Link struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	URL       string `json:"url"`
	Title     string `json:"title,omitempty"`
	Clicks    int    `json:"clicks"`
	CreatedAt string `json:"created_at"`
	LastClick string `json:"last_click,omitempty"`
}
type ClickLog struct {
	ID        string `json:"id"`
	LinkID    string `json:"link_id"`
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referer   string `json:"referer,omitempty"`
	CreatedAt string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "crossroads.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS links(id TEXT PRIMARY KEY,slug TEXT UNIQUE NOT NULL,url TEXT NOT NULL,title TEXT DEFAULT '',clicks INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')),last_click TEXT DEFAULT '')`,
		`CREATE TABLE IF NOT EXISTS click_log(id TEXT PRIMARY KEY,link_id TEXT NOT NULL,ip TEXT DEFAULT '',user_agent TEXT DEFAULT '',referer TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`,
		`CREATE INDEX IF NOT EXISTS idx_clicks_link ON click_log(link_id)`,
	} {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(resource TEXT NOT NULL,record_id TEXT NOT NULL,data TEXT NOT NULL DEFAULT '{}',PRIMARY KEY(resource, record_id))`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func genSlug() string      { b := make([]byte, 4); rand.Read(b); return hex.EncodeToString(b)[:6] }
func (d *DB) Create(l *Link) error {
	l.ID = genID()
	l.CreatedAt = now()
	if l.Slug == "" {
		l.Slug = genSlug()
	}
	_, err := d.db.Exec(`INSERT INTO links(id,slug,url,title,created_at)VALUES(?,?,?,?,?)`, l.ID, l.Slug, l.URL, l.Title, l.CreatedAt)
	return err
}
func (d *DB) GetBySlug(slug string) *Link {
	var l Link
	if d.db.QueryRow(`SELECT id,slug,url,title,clicks,created_at,last_click FROM links WHERE slug=?`, slug).Scan(&l.ID, &l.Slug, &l.URL, &l.Title, &l.Clicks, &l.CreatedAt, &l.LastClick) != nil {
		return nil
	}
	return &l
}
func (d *DB) GetByID(id string) *Link {
	var l Link
	if d.db.QueryRow(`SELECT id,slug,url,title,clicks,created_at,last_click FROM links WHERE id=?`, id).Scan(&l.ID, &l.Slug, &l.URL, &l.Title, &l.Clicks, &l.CreatedAt, &l.LastClick) != nil {
		return nil
	}
	return &l
}
func (d *DB) List() []Link {
	rows, _ := d.db.Query(`SELECT id,slug,url,title,clicks,created_at,last_click FROM links ORDER BY created_at DESC`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Link
	for rows.Next() {
		var l Link
		rows.Scan(&l.ID, &l.Slug, &l.URL, &l.Title, &l.Clicks, &l.CreatedAt, &l.LastClick)
		o = append(o, l)
	}
	return o
}
func (d *DB) Delete(id string) error {
	d.db.Exec(`DELETE FROM click_log WHERE link_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM links WHERE id=?`, id)
	return err
}
func (d *DB) RecordClick(slug, ip, ua, ref string) *Link {
	l := d.GetBySlug(slug)
	if l == nil {
		return nil
	}
	t := now()
	d.db.Exec(`UPDATE links SET clicks=clicks+1,last_click=? WHERE id=?`, t, l.ID)
	d.db.Exec(`INSERT INTO click_log(id,link_id,ip,user_agent,referer,created_at)VALUES(?,?,?,?,?,?)`, genID(), l.ID, ip, ua, ref, t)
	return d.GetBySlug(slug)
}
func (d *DB) ClickHistory(linkID string, limit int) []ClickLog {
	if limit <= 0 {
		limit = 50
	}
	rows, _ := d.db.Query(`SELECT id,link_id,ip,user_agent,referer,created_at FROM click_log WHERE link_id=? ORDER BY created_at DESC LIMIT ?`, linkID, limit)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []ClickLog
	for rows.Next() {
		var c ClickLog
		rows.Scan(&c.ID, &c.LinkID, &c.IP, &c.UserAgent, &c.Referer, &c.CreatedAt)
		o = append(o, c)
	}
	return o
}

type Stats struct {
	Links       int `json:"links"`
	TotalClicks int `json:"total_clicks"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM links`).Scan(&s.Links)
	d.db.QueryRow(`SELECT COALESCE(SUM(clicks),0) FROM links`).Scan(&s.TotalClicks)
	return s
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
