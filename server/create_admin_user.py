#!/usr/bin/env python3

import argparse
import sqlite3
import bcrypt
import sys
import os

def main():
    parser = argparse.ArgumentParser(description='Add a new admin user to the database.')
    parser.add_argument('db_path', help='Path to the SQLite database file (e.g., lbd_game.db)')
    parser.add_argument('username', help='Username for the new admin user')
    parser.add_argument('password', help='Password for the new admin user')

    args = parser.parse_args()

    db_path = args.db_path
    username = args.username
    password = args.password

    # Check if the database file exists
    if not os.path.isfile(db_path):
        print(f"Error: Database file '{db_path}' does not exist.")
        sys.exit(1)

    # Connect to the SQLite database
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # Create the User table if it doesn't exist
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            is_approved BOOLEAN DEFAULT FALSE,
            is_eliminated BOOLEAN DEFAULT FALSE,
            is_admin BOOLEAN DEFAULT FALSE
        )
    ''')

    # Check if the username already exists
    cursor.execute('SELECT id FROM users WHERE username = ?', (username,))
    if cursor.fetchone():
        print(f"Error: Username '{username}' already exists.")
        conn.close()
        sys.exit(1)

    # Hash the password using bcrypt
    password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
    password_hash_str = password_hash.decode('utf-8')

    # Insert the new admin user into the database
    cursor.execute('''
        INSERT INTO users (username, password, is_approved, is_eliminated, is_admin)
        VALUES (?, ?, ?, ?, ?)
    ''', (username, password_hash_str, True, False, True))

    conn.commit()
    conn.close()

    print(f"Admin user '{username}' added successfully.")

if __name__ == '__main__':
    main()
