phone-yml:
	kubectl create secret generic phone-numbers -n quotes --from-literal=numbers=$(N) --dry-run=client -o yaml \
		| kubeseal --format=yaml - \
		> ./ops/base/phoneSecret.yml

apply-phone: phone-yml
	kubectl apply -f ./ops/base/phoneSecret.yml
