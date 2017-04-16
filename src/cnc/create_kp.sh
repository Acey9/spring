read -p "Enter your domain [www.example.com]: " DOMAIN

echo "Create server key..."

openssl genrsa -des3 -out $DOMAIN.key 1024

echo "Create server certificate signing request..."

SUBJECT="/C=CN/ST=SP/L=spring/O=spring/OU=spring/CN=$DOMAIN"

openssl req -new -subj $SUBJECT -key $DOMAIN.key -out $DOMAIN.csr

echo "Remove password..."

mv $DOMAIN.key $DOMAIN.origin.key
openssl rsa -in $DOMAIN.origin.key -out $DOMAIN.key

echo "Sign SSL certificate..."

openssl x509 -req -days 3650 -in $DOMAIN.csr -signkey $DOMAIN.key -out $DOMAIN.crt
