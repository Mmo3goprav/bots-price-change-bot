# 📈 CryptoChartBot – Telegram Bot for Monitoring Currency Pairs  

**CryptoChartBot** is a Telegram bot that helps you track cryptocurrency price changes in real time.

## 🚀 Bot Features
✅ **`/start`** – Starts the bot and displays the available commands.  
✅ **"➕ Add Chart"** – Adds a currency pair for tracking.  
✅ **"❌ Remove Chart"** – Removes a currency pair from monitoring.  

## 🔔 How Does the Bot Work?
- Once a currency pair is added, the bot sends **automatic updates** in the following format:
```
Symbol: BTCUSD
Current price: 85470.000000
Price change: 1149.000000
Price change percentage: 0.000000
```


- The data updates in real time and is sent directly to the user.

## 📌 How to Use the Bot?
1️⃣ **Start the bot** by using the `/start` command.  
2️⃣ **Add a currency pair** using the "Add Chart" button.  
3️⃣ **Receive real-time updates** about price changes in your chat.  
4️⃣ **Remove a pair** when you no longer need it by clicking "Remove Chart".  

## 🛠 How to Run the Bot?
1️⃣ **Clone the repository:**
 ```sh
 git clone https://github.com/Mmo3goprav/bots-price-change-bot.git
 cd YOUR_REPOSITORY
```
2️⃣ **Set up your bot token:**

- Create a `.env` file and add your Telegram bot token:
  ```sh
  TOKEN=your_token_here
  ```
- Or manually set the token in your code.

3️⃣ **Run the bot:**
```sh
go run main.go
```

⚡ **CryptoChartBot** – a fast and convenient way to track cryptocurrency prices directly in Telegram!
