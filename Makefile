combine:
	find . -type f -name '*.go' -print0 | sort -z | xargs -0 cat > combined.go
