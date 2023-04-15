include .env
export

phone-yml:
	kubectl create secret generic phone-numbers -n quotes --from-literal=numbers=$(N) --dry-run=client -o yaml \
		| kubeseal --format=yaml - \
		> ./ops/base/phoneSecret.yml

apply-phone: phone-yml
	kubectl apply -f ./ops/base/phoneSecret.yml

build:
	rm -rf dist/*
	go mod download
	mkdir -p ./dist
	go build -o ./dist/quoteCats

deploy:
	mkdir -p ./vanilla_deploy/tmp

	cat ./vanilla_deploy/env.template \
		| sed "s/TWILIO_AUTH_VAR/${TWILIO_AUTH}/g" \
		| sed "s/PHONE_NUMBERS_VAR/${PHONE_NUMBERS}/g" \
		| sed "s/QUOTE_API_KEY_VAR/${QUOTE_API_KEY}/g" \
		| sed "s/API_NINJA_KEY_VAR/${API_NINJA_KEY}/g" \
		| sed "s/CAT_API_KEY_VAR/${CAT_API_KEY}/g" > ./vanilla_deploy/tmp/.env

	cat ./vanilla_deploy/loki-config.yaml.template \
		| sed "s/GRAFANA_LOKI_API_KEY_VAR/${GRAFANA_LOKI_API_KEY}/g" > ./vanilla_deploy/tmp/promtail-config.yaml
	
	scp ./vanilla_deploy/startup.sh ${REMOTE_MACHINE_IP}:${USER_DIR}/startup.sh
	scp ./vanilla_deploy/tmp/.env ${REMOTE_MACHINE_IP}:${USER_DIR}/.env
	scp ./vanilla_deploy/tmp/promtail-config.yaml ${REMOTE_MACHINE_IP}:${USER_DIR}/promtail-config.yaml
	scp ./vanilla_deploy/promtail.service ${REMOTE_MACHINE_IP}:${USER_DIR}/promtail.service

	ssh ${REMOTE_MACHINE_IP} ${USER_DIR}/startup.sh

	rm -rf ./vanilla_deploy/tmp

remote:
	sed -i '' 's/REMOTE_MACHINE_IP=.*/REMOTE_MACHINE_IP=$(shell gcloud compute instances describe "quote-sender" --billing-project="quote-sender-381016" --format="json" --zone="us-west1-a" | jq ".networkInterfaces[].accessConfigs[].natIP")/' .env

ssh:
	ssh $(shell gcloud compute instances describe 'quote-sender' --billing-project='quote-sender-381016' --format='json' --zone='us-west1-a' | jq '.networkInterfaces[].accessConfigs[].natIP')
