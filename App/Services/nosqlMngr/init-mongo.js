// Initialize MongoDB collections and insert sample data
db = db.getSiblingDB("tweetsdb");

// Create tweets collection
db.createCollection("tweets");

// Insert sample tweets
db.tweets.insertMany([
  { owner: 1, tweet: "Hello Tweet1 Hello Tweet1 Hello Tweet1", created_at: new Date() },
  { owner: 2, tweet: "Hello Tweet2 Hello Tweet2", created_at: new Date() }
]);