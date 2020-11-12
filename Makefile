.PHONY = apply all clean

# List all subdirectories that we want to run "go build" on
images := receiver sender

targets := $(foreach dir,$(images),$(dir)/$(dir))

all: $(targets)

apply-%: all
	${GOPATH}/bin/ko apply -f config/$(subst apply-,,$@)/

%: %.go
	cd $(dir $@); go build

clean:
	rm -f $(targets)
