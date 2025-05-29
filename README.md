<h1>Golang Order Matching System</h1>
<h2>Postman Collection file Also provided in Repository for all API EndPoint Testing(Order_matching_system.postman_collection.json)</h2>
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

<h2>API EndPoints</h2>
<h4>Base URL: http://localhost:8080/api/v1</h3>
<pre><code>1. Place Order
http 
POST /orders
Content-Type: application/json

{
    "symbol": "BTCUSD",
    "side": "buy",
    "type": "limit",
    "price": "50000.00",
    "quantity": "0.1"
}</code></pre>
<pre><code>2. Cancel Order
http
DELETE /orders/{orderId}</code></pre>
<pre><code>3. Get Order Status
http
GET /orders/{orderId}</code></pre>

<pre><code> 4. Get Order Book
http
GET /orderbook?symbol=BTCUSD
</code></pre>

<pre><code>5. GetTrades
http
GET /trades?symbol=BTCUSD&limit=50</code></pre>
<h2>Results</h2>
<pre><code> <h3>PlaceOrder </h3>
<img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/Place_BUY_LIMIT_ORDER.PNG"/>
 <img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/place_sell_limit_order.PNG"/></code></pre>
<pre><code> <h3>GetOrderBook </h3>
<img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/GET_ORDER_BOOK.PNG"/></code></pre>
<pre><code> <h3>GetTrade </h3>
<img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/get_trade.PNG"/></code></pre>
<pre><code> <h3>Edge Cases:</h3>
<img src="https://github.com/spee-dev/GOLANG-ORDER-MATCHING-SYSTEM/blob/main/edge_case.PNG"/></code></pre>

