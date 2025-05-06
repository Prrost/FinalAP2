module github.com/Prrost/FinalAP2/book-service

go 1.23.0

toolchain go1.24.2

//replace github.com/Prrost/protoFinalAP2 => ../protoFinalAP2

require (
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.28
)

require github.com/Prrost/protoFinalAP2 v0.0.0-20250505065838-82b7b58ea42e // indirect
