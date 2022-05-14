# fake-github-webhook

# What is?
This is an application that simulates Github webhooks.
You choose a folder containing JSON files with events payload and then send it to a given target.
Each JSON file should contain a sequence of events.

A Sequence is an array of Github Webhook events. The application execute events in order within a Sequence. 

# In which scenarios I can use this thing?
- You can use it to test Continous Integration tools on localhost. This is good because you can test your tool and put it to run on internet only when it's ready.
- Also, you can use it to avoid useless commits on a test repository.

# How I use this thing?
    # start your CI tool listening on 127.0.0.1:8080
    # then test it triggering a payload to 127.0.0.1:8080
    $ fake-github-webhook -host 127.0.0.1:8080 -data-dir data -interval 3s

