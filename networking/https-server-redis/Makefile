export CERT_DIR:=$(HOME)/Documents/.localhost_cert
dc:
	echo "Executing  docker-compose up -d --build"
	docker-compose up -d --build
clean:
	echo "Executing docker-compose down --rmi 'local'"
	docker-compose down --rmi 'local'

