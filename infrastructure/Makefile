setup:
	terraform remote -address="$TF_STATE_SERVER/playground" -backend=HTTP

pull:
	terraform pull

push:
	terraform push

plan:
	terraform plan -out plan.tfplan

destroy:
	terraform plan -destroy -out plan.tfplan

apply:
	terraform apply plan.tfplan
