.PHONY: mocks machine

# Generation Variables
MOCKERY_CMD = mockery
MACHINE_CMD = go run cmd/external-device/main.go

mocks:
	@echo "Generating mocks with mockery..."
	$(MOCKERY_CMD)

## Application Execution Rules
## ---------------------------

machine:
	@echo "Starting external-device machine..."
	$(MACHINE_CMD)