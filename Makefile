STACK = go-trace
PROFILE = mine

.PHONY: build deploy validate destroy

deploy:
	sam build
	AWS_PROFILE=$(PROFILE) sam deploy --stack-name $(STACK) \
	--resolve-s3 --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM

validate:
	AWS_PROFILE=$(PROFILE) aws cloudformation validate-template --template-body file://template.yml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

clean:
	rm -rf main gin-bin

