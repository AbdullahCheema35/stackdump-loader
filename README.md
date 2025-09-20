# stackdump-loader

`stackdump-loader` is a Go-based toolchain for **processing and importing Stack Exchange data dumps** into PostgreSQL.  
It converts the large XML files (from [Stack Exchange Data Dump](https://archive.org/details/stackexchange)) into **chunked CSV files** (10M rows each), and then automatically streams them into PostgreSQL using `\copy`.  

This makes it easier to handle massive datasets like Stack Overflow, Super User, Server Fault, etc. without running into memory or file-size issues.

---

## 📦 Features
- Stream-parse huge Stack Exchange XML files (no need to load entire XML in memory).
- Convert to CSV with proper headers.
- Split output into multiple CSVs (default: 10M rows each).
- Automatically run `make posts filename=<csv-file>` to import into PostgreSQL.
- Supports multiple Stack Exchange dump files:
  - `Users.xml`
  - `Tags.xml`
  - `Badges.xml`
  - `Votes.xml`
  - `Comments.xml`
- Deletes CSV after successful import (to save disk space).

---

## 🗂 Stack Exchange Data Dumps
Official dumps are published here:

👉 [Stack Exchange Data Dump @ Archive.org](https://archive.org/details/stackexchange)

Each dump is a `.7z` archive per site (e.g., `stackoverflow.com-Posts.7z`).  
You can extract these files with [7-Zip](https://www.7-zip.org/) or `p7zip` on Linux:

```bash
7z x stackoverflow.com-Users.7z
