<h1>Golang Order Matching System</h1>
<h2>Postman Collection file Also provided in Repository for API Tesrting</h2>
<h2>Project Setup Instructions</h2>

<ol>
  <li><strong>Install Prerequisites</strong>
    <ul>
      <li>Go (v1.18+): <a href="https://golang.org/dl/">https://golang.org/dl/</a></li>
      <li>MySQL (v5.7+): <a href="https://dev.mysql.com/downloads/">https://dev.mysql.com/downloads/</a></li>
      <li>Git (optional): <a href="https://git-scm.com/">https://git-scm.com/</a></li>
    </ul>
  </li>

  <li><strong>Clone the Repository</strong>
    <pre><code>git clone https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM.git
cd GOLANG-ORDER-MATCHING-SYSTEM</code></pre>
  </li>

  <li><strong>Environment Configuration</strong>
    <p>Edit <code>.env</code>:</p>
    <pre><code>PORT=8080
DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ordermatching</code></pre>
  </li>

  <li><strong>Start MySQL & Create Database</strong>
    <p>Login to MySQL and run:</p>
    <pre><code>CREATE DATABASE ordermatching;</code></pre>
  </li>

  <li><strong>Install Go Dependencies</strong>
    <pre><code>go mod tidy</code></pre>
  </li>

  <li><strong>Run the Application</strong>
    <pre><code>go run cmd/server/main.go</code></pre>
    <p>This will:</p>
    <ul>
      <li>Load environment variables</li>
      <li>Connect to MySQL</li>
      <li>Auto-run SQL migrations</li>
      <li>Start HTTP server at <code>localhost:8080</code></li>
    </ul>
  </li>
</ol>

<h2>API Usage (Sample Calls)</h2>

<h3>Place an Order</h3>
<pre><code>POST http://localhost:8080/api/v1/orders \
-H "Content-Type: application/json" \
-d '{
  "symbol": "BTCUSD",
  "side": "buy",
  "type": "limit",
  "price": "28000",
  "quantity": "1"
}'</code></pre>

<h3>Cancel an Order</h3>
<pre><code>DELETE http://localhost:8080/api/v1/orders/{orderID}</code></pre>

<h3>Get Order Book</h3>
<pre><code> http://localhost:8080/api/v1/orderbook?symbol=BTCUSD</code></pre>

<h3>Get Trades</h3>
<pre><code> http://localhost:8080/api/v1/trades?symbol=BTCUSD</code></pre>


<h2>OutPut</h2>
<img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/Place_BUY_LIMIT_ORDER.PNG"/>

