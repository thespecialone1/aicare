// index.js
const { Client, LocalAuth } = require('whatsapp-web.js');
const axios = require('axios');
const puppeteer = require('puppeteer');
const qrcode = require('qrcode-terminal');

// 1) Configure these:
const API_URL = 'http://localhost:8080/question';  // your Go backend
const PREFIX_QA   = '#aicare ';
const PREFIX_DIAG = '#aicare-diag ';

// 2) Initialize WhatsApp client
const client = new Client({
  authStrategy: new LocalAuth(),
  puppeteer: {  
    // use the full puppeteer package you installed
    executablePath: puppeteer.executablePath(),
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
    ],
  },
});

// 3) On QR code (first-time only)
client.on('qr', qr => {
  console.log('Scan this QR code with your WhatsApp app:');
  qrcode.generate(qr, { small: true });  // renders the QR in your terminal
});

// 4) On ready
client.on('ready', () => {
  console.log('WhatsApp client ready');
});

// 5) On incoming message
client.on('message', async msg => {
  const body = msg.body || '';
  
  // Determine mode and extract question
  let mode = null
  let question = null
  
  if (body.toLowerCase().startsWith(PREFIX_DIAG)) {
    mode = 'diag';
    question = body.slice(PREFIX_DIAG.length).trim()
  } else if (body.toLowerCase().startsWith(PREFIX_QA)) {
    mode = 'answer';
    question = body.slice(PREFIX_QA.length).trim();
  } else {
    return; // ignore messages without our prefixes
  }

  console.log(`Mode: ${mode}, Question: ${question}`);
  
try {
    // 6) Call your backend API with both question and mode
    const res = await axios.post(API_URL, {
      question,
      mode
    }, {
      headers: { 'Content-Type': 'application/json' }
    });

    const answer = res.data.answer;
    console.log(`Replying: ${answer}`);

    // 7) Reply on WhatsApp
    await msg.reply(answer);
  } catch (err) {
    console.error('Error calling API:', err.message);
    await msg.reply('❌ Sorry, I couldn’t process your question right now.');
  }
});
// 9) Start the client
client.initialize();
