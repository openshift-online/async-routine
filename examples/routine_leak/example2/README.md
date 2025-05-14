# Routine Leak Detection with AsyncRoutineManager
This demonstration illustrates how the `AsyncRoutineManager` help identifying routine leaks and pinpoint their source.

Each folder, named step1 to stepN, adds a progressive integration of the `AsyncRoutineManager`, starting from the 
naked code of `step1` to the full integration in the last step.

The application we are going to use for this demonstration is very simple:

1. The `main` function repeatedly starts cycles where it processes multiple websites.
   For each cycle, randomly selects 10 website URLs from a predefined list.
   For each selected URL, asynchronously invokes the `getWebsiteResponseSize`
2. The `getWebsiteResponseSize` asynchronously calls the `getResponseSize` and do some fun stuff with the result
3. The `getResponseSize` contacts the site and returns the response size
