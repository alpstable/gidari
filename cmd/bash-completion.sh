#!/bin/bash

# __bashcompembed_custom_func is a function that is called by bash-completion to get a list of completions used by the
# legacy bash autocomplete system. This function is irrelevant to the new bash autocomplete system and should do
# nothing.
__bashcompembed_custom_func() {
	ehco "This function is irrelevant to the new bash autocomplete system and should return nothing."
}
