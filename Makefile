.PHONY: mocks machine inspect

# Generation Variables
MOCKERY_CMD = mockery
MACHINE_CMD = go run cmd/external-device/main.go
BADGER_INSPECT_CMD = go run cmd/inspect/main.go

mocks:
	@echo "Generating mocks with mockery..."
	$(MOCKERY_CMD)

## Application Execution Rules
## ---------------------------

machine:
	@echo "Starting external-device machine..."
	$(MACHINE_CMD)

inspect:
	$(BADGER_INSPECT_CMD)