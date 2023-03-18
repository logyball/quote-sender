include .env
export

phone-yml:
	kubectl create secret generic phone-numbers -n quotes --from-literal=numbers=$(N) --dry-run=client -o yaml \
		| kubeseal --format=yaml - \
		> ./ops/base/phoneSecret.yml

apply-phone: phone-yml
	kubectl apply -f ./ops/base/phoneSecret.yml

build:
	go mod download
	mkdir ./dist
	go build -o ./dist/quoteCats

deploy:
	mkdir -p ./vanilla_deploy/tmp
	cat ./vanilla_deploy/env.template \
		| sed "s/TWILIO_AUTH_VAR/${TWILIO_AUTH}/g" \
		| sed "s/PHONE_NUMBERS_VAR/${PHONE_NUMBERS}/g" \
		| sed "s/CAT_API_KEY_VAR/${CAT_API_KEY}/g" > ./vanilla_deploy/tmp/.env
	cat ./vanilla_deploy/tmp/.env
	
	scp ./vanilla_deploy/startup.sh ${REMOTE_MACHINE_IP}:${USER_DIR}/startup.sh
	scp ./vanilla_deploy/tmp/.env ${REMOTE_MACHINE_IP}:${USER_DIR}/.env

	ssh ${REMOTE_MACHINE_IP} ${USER_DIR}/startup.sh

	rm -rf ./vanilla_deploy/tmp