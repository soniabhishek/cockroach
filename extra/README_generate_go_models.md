
WELCOME TO pg2go (generate_go_models)
====================================

## Whats the point?

It auto-generates the go models through the database tables. It adds db, json & bson tags also.
Also you can add your custom rules. Go look at pg2go.sql

### How to run?

# Windows/Ubuntu/Linux

- Start git bash/terminal & make sure you have psql & goimports in your path
- run it from angel's root :
` source extra/generate_go_models.sh <DBNAME> <DBUSERNAME> `

e.g.
` source extra/generate_go_models.sh playment_local postgres `

- The above command will update app/models/db_struct.go
- It will also add CREATE FUNCTION & DROP FUNCTION texts at the end & bottom of file.
    - Remove that manaully 
    - rerun goimports on that file: 
         `  goimports -w "$OUT" || gofmt -w "$OUT"  `
