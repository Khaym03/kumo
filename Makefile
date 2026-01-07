.PHONY: mocks test machine inspect

# Generation Variables
MOCKERY_CMD = mockery
MACHINE_CMD = go run cmd/external-device/main.go
BADGER_INSPECT_CMD = go run cmd/inspect/main.go
TEST_CMD = go test ./... -v

mocks:
	@echo "Generating mocks with mockery..."
	$(MOCKERY_CMD)

test:
	$(TEST_CMD)

## Application Execution Rules
## ---------------------------

machine:
	@echo "Starting external-device machine..."
	$(MACHINE_CMD)

inspect:
	$(BADGER_INSPECT_CMD)
