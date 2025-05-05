# Receipt Processor Challenge - JavaScript 

My default answer will be in Javascript since that is the language I am the most familiar with but due to the nature of this opputinity I did take the time to go the extra mile and learn a bit of Golang as well. There will be a folder with what my answer would look like in Go as I did take my time to learn the basics of the language and utilize the time I did have to think like someone who was already on the team so I did take the time to Familarize myself with Go and Docker!

## Prerequisites
- Docker should be installed on your machine.
- The application is built with Node.js and Express.

## Steps to Run

1. Clone this repository or download the files.

2. Build the Docker image:
   ```bash
   docker build -t receipt-processor-challenge .

3. Run The Docker Container:
   docker run -p 3000:3000 receipt-processor-challenge

4. The application will now be running on http://localhost:3000. You can make requests to:
   
   POST /receipts/process to process a receipt.

   GET /receipts/points/:id to get points for a receipt by ID.

   GET / to view all receipts.

