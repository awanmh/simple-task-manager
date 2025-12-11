package domain

import "errors"

// ErrNotFound adalah error standar ketika data tidak ditemukan di DB
var ErrNotFound = errors.New("record not found")

// ErrConflict bisa kita tambahkan juga untuk kasus duplikat data (misal email sama)
var ErrConflict = errors.New("record already exists")