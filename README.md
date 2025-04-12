# AkaveAI

![AkaveAI Front Image](/akaveai-front.png)

## Introduction  
AkaveAI is a decentralized marketplace that enables developers to create, publish, and monetize AI Agents securely. Built using Filecoin and Akave, it allows users to discover, purchase, and deploy autonomous AI agents for various tasks—ensuring transparency, ethical AI use, and fair rewards for creators.

This platform bridges the gap between AI agent developers and users, fostering a decentralized, privacy-first AI ecosystem.

## Features

**For AI Developers:**  
- Upload AI agent metadata securely to Filecoin using Akave.  
- Set a price for access and receive fair compensation.  
- Retain credit and ownership rights over your agents.

**For Users:**  
- Browse and purchase AI agents with cryptocurrency.  
- Instantly access agents after secure blockchain-based payment.  
- Ensure compliance with ethical AI sourcing and attribution.

**Security & Transparency:**  
- Uses smart contracts for automated payments and access control.  
- Ensures immutable and verifiable ownership records.  
- Decentralized storage prevents unauthorized manipulation.

## How It Works

**Step 1: Agent Upload**  
Developers upload agent configuration files and metadata to Filecoin via Akave.  
The platform generates a unique agent ID and stores metadata.

**Step 2: Agent Listing**  
The agent appears in the marketplace with details (name, task, description, price, etc.).  
Users can explore and filter agents based on capabilities.

**Step 3: Purchase & Payment**  
Users select an agent and pay using cryptocurrency.  
A smart contract verifies the transaction and releases access.

**Step 4: Access & Monetization**  
Buyers receive a secure endpoint or config link.  
The developer receives payment (e.g., 85% revenue share).  
The platform takes a small fee for maintenance.

## Smart Contract Details

**Agent Management**
- `listAgent(bucketName, fileName, price)`: Adds an agent to the marketplace.  
- `buyAgent(agentID)`: Purchases agent access and distributes payment.  
- `getAgent(agentID)`: Retrieves agent metadata.  
- `getUserPurchases(userAddress)`: Returns a list of agents purchased by a user.

**Earnings & Statistics**
- `getAgentSales(ownerAddress)`: Returns the number of sales for a developer.  
- `getDeveloperEarnings(ownerAddress)`: Returns total earnings from agent sales.

## Usage Instructions

**Deployment**
- Deploy the smart contract on Ethereum or any EVM-compatible chain.  
- Connect the contract to the frontend using Web3.  
- Integrate Filecoin via Akave for decentralized storage.

**Uploading Agents**
- Navigate to the Upload section.  
- Enter agent name, description, price, and select the config file.  
- Submit to store data on Filecoin and list it on the marketplace.

**Buying Agents**
- Browse available AI agents.  
- Select an agent and confirm payment using crypto.  
- Receive access details after successful payment.

## Events

- `AgentListed(uint256 id, string bucketName, string fileName, uint256 price, address owner)`: An agent was listed.  
- `AgentPurchased(uint256 id, address buyer)`: An agent was purchased.

## Security Considerations

- Smart contracts are designed to prevent unauthorized access or withdrawals.  
- Payments are automated via blockchain to ensure fair transactions.  
- Agent metadata remains immutable on decentralized storage.  
- Regular audits are recommended before mainnet deployment.

## Tech Stack

- **Blockchain:** Ethereum / EVM-compatible chains  
- **Smart Contracts:** Solidity  
- **Storage:** Filecoin (via Akave API)  
- **Frontend:** React.js + Web3.js  
- **Backend:** Node.js + Express

## Future Enhancements

- Add reputation-based pricing and agent ratings.  
- Introduce NFT-based licensing for AI agents.  
- Support multi-chain deployment.  
- Enable auto-updating agent capabilities.  
- Add usage-based billing models.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

## Support

For issues and support:  
- Open an issue on GitHub  
- Contact us via email at anishgajbhare2000@gmail.com

## Links

- [GitHub Repository: AkaveAI](https://github.com/anishgajbhare/akaveai)  
- [Demo Video: Watch Here](https://www.youtube.com/watch?v=your-demo-video)  
- [Project Website: Visit Here](https://akaveai.xyz)
