#!/usr/bin/env python3

import argparse
import sqlite3
import bcrypt
import sys
import os

def main():
    parser = argparse.ArgumentParser(description='Change the password of an existing user.')
    parser.add_argument('db_path', help='Path to the SQLite database file (e.g., lbd_game.db)')
    parser.add_argument('username', help='Username of the user')
    parser.add_argument('new_password', help='New password for the user')

    args = parser.parse_args()

    db_path = args.db_path
    username = args.username
    new_password = args.new_password

    # Check if the database file exists
    if not os.path.isfile(db_path):
        print(f"Error: Database file '{db_path}' does not exist.")
        sys.exit(1)

    # Connect to the SQLite database
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    # Check if the users table exists
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

    # Check if the user exists
    cursor.execute('SELECT id FROM users WHERE username = ?', (username,))
    result = cursor.fetchone()
    if not result:
        print(f"Error: Username '{username}' does not exist.")
        conn.close()
        sys.exit(1)

    user_id = result[0]

    # Hash the new password using bcrypt
    new_password_hash = bcrypt.hashpw(new_password.encode('utf-8'), bcrypt.gensalt())
    new_password_hash_str = new_password_hash.decode('utf-8')

    # Update the user's password in the database
    cursor.execute('''
        UPDATE users
        SET password = ?
        WHERE id = ?
    ''', (new_password_hash_str, user_id))

    conn.commit()
    conn.close()

    print(f"Password for user '{username}' has been updated successfully.")

if __name__ == '__main__':
    main()
