package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
)

// DB 定义数据库结构体
type DB struct {
	conn *sql.DB
}

// NewDB 创建并初始化 SQLite 数据库连接
func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	// 创建必要的表
	if err := createTables(conn); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// createTables 创建域名和证书检查历史记录的表
func createTables(conn *sql.DB) error {
	createDomainTable := `
	CREATE TABLE IF NOT EXISTS domains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		notify_before_days INTEGER NOT NULL
	);`

	createCheckHistoryTable := `
	CREATE TABLE IF NOT EXISTS check_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain_id INTEGER NOT NULL,
		checked_at DATETIME NOT NULL,
		expiration_date DATETIME NOT NULL,
		days_left INTEGER NOT NULL,
		FOREIGN KEY(domain_id) REFERENCES domains(id)
	);`

	_, err := conn.Exec(createDomainTable)
	if err != nil {
		return fmt.Errorf("failed to create domains table: %v", err)
	}

	_, err = conn.Exec(createCheckHistoryTable)
	if err != nil {
		return fmt.Errorf("failed to create check_history table: %v", err)
	}

	return nil
}

// AddDomain 添加一个新的域名到数据库
func (db *DB) AddDomain(name string, notifyBeforeDays int) error {
	query := `INSERT OR IGNORE INTO domains (name, notify_before_days) VALUES (?, ?)`
	_, err := db.conn.Exec(query, name, notifyBeforeDays)
	if err != nil {
		return fmt.Errorf("failed to insert domain: %v", err)
	}
	return nil
}

// GetDomains 获取所有的域名列表
func (db *DB) GetDomains() ([]map[string]interface{}, error) {
	query := `SELECT id, name, notify_before_days FROM domains`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query domains: %v", err)
	}
	defer rows.Close()

	var domains []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var notifyBeforeDays int
		if err := rows.Scan(&id, &name, &notifyBeforeDays); err != nil {
			return nil, fmt.Errorf("failed to scan domain: %v", err)
		}
		domain := map[string]interface{}{
			"id":                 id,
			"name":               name,
			"notify_before_days": notifyBeforeDays,
		}
		domains = append(domains, domain)
	}
	return domains, nil
}

// AddCheckHistory 添加证书检查记录
func (db *DB) AddCheckHistory(domainID int, expirationDate time.Time, daysLeft int) error {
	query := `INSERT INTO check_history (domain_id, checked_at, expiration_date, days_left) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, domainID, time.Now(), expirationDate, daysLeft)
	if err != nil {
		return fmt.Errorf("failed to insert check history: %v", err)
	}
	return nil
}

// GetCheckHistory 获取某个域名的证书检查历史记录
func (db *DB) GetCheckHistory(domainID int) ([]map[string]interface{}, error) {
	query := `SELECT checked_at, expiration_date, days_left FROM check_history WHERE domain_id = ? ORDER BY checked_at DESC`
	rows, err := db.conn.Query(query, domainID)
	if err != nil {
		return nil, fmt.Errorf("failed to query check history: %v", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var checkedAt time.Time
		var expirationDate time.Time
		var daysLeft int
		if err := rows.Scan(&checkedAt, &expirationDate, &daysLeft); err != nil {
			return nil, fmt.Errorf("failed to scan check history: %v", err)
		}
		record := map[string]interface{}{
			"checked_at":      checkedAt,
			"expiration_date": expirationDate,
			"days_left":       daysLeft,
		}
		history = append(history, record)
	}
	return history, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}
	return nil
}
