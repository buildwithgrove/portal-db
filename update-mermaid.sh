#!/bin/bash

mermerd -c "postgres://postgres:pgpassword@localhost:5432/postgres?sslmode=disable" -s public --useAllTables # pragma: allowlist secret
sed -i '/# Current Schema/,$d' README.md
echo -e "# Current Schema\n\n\`\`\`mermaid" >> README.md
cat result.mmd >> README.md
echo -e "\`\`\`\n" >> README.md
rm result.mmd
