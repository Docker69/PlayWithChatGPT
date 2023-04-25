set -e

mongosh <<EOF

print("============= Started INIT script =======================");

db = db.getSiblingDB('$MONGO_INITDB_DATABASE')

db.createUser({
  user: '$MONGO_USER_NAME',
  pwd: '$MONGO_USER_PASSWORD',
  roles: [{ role: 'readWrite', db: '$MONGO_INITDB_DATABASE' }],
});

db.createCollection('chats')
db.createCollection('configs')
db.createCollection('humans')

print("============= Ended INIT script =======================");

EOF