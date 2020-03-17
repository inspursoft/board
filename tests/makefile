ifeq ($(DIR),)
	CLONEDIR=empty
else
	CLONEDIR=$(DIR)
endif

e:
	@echo "no parameter"
all:
	@echo $(CLONEDIR)
	@/bin/bash pull.sh $(CLONEDIR)
	@/bin/bash run.sh 
run:
	@/bin/bash run.sh
clone:
	@/bin/bash pull.sh $(CLONEDIR)
