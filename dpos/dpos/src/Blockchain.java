import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.Random;
import java.util.UUID;

public class Blockchain {
    // 区块链列表
    private static List<Blockchain> blockchainList = new ArrayList<>();
    // 区块链的难度
    private int difficulty;
    // 投票列表
    private List<Vote> voteList = new ArrayList<>();
    // 节点列表
    private List<Node> nodeList = new ArrayList<>();
    // 区块列表
    private List<Block> blockList = new ArrayList<>();
    //交易列表
    private List<Transaction> transaction = new ArrayList<>();
    
    public static class Block {
        private int index;
        private long timestamp;
        private List<Transaction> transactionList;
        private List<Vote> voteList;
        private String previousHash;
        private String hash;
        private int nonce;
    
        public Block(int index, long timestamp, List<Transaction> transactionList, List<Vote> voteList, String previousHash,
                int nonce) {
            this.index = index; 
            this.timestamp = timestamp;
            this.transactionList = transactionList;
            this.voteList = voteList;
            this.previousHash = previousHash;
            this.nonce = nonce;
            //getHash
            this.hash = HashUtils.sha256(String.valueOf(timestamp));
        }
    
    }

    public void createGenesisBlock() {
        List<Transaction> transactions = new ArrayList<>();
        List<Vote> votes = new ArrayList<>();
        String previousHash = "";
        int nonce = 0;
        Block genesisBlock = new Block(0, System.currentTimeMillis(), transactions, votes, previousHash, nonce);
        blockList.add(genesisBlock);
        System.out.println("genesisBlock: " + genesisBlock.hash);
    }

    public void addNode(Node node) {
        nodeList.add(node);
    }
    public void addVote(Vote vote) {
        voteList.add(vote);
    }

    public List<Node> sortNodesByVoteCount() {
        Collections.sort(nodeList, new Comparator<Node>() {
            @Override
            public int compare(Node node1, Node node2) {
                return Integer.compare(node2.getVoteCount(), node1.getVoteCount());
            }
        });
        return nodeList;
    }

    public void addBlock(Block block) {
        blockList.add(block);
    }

    public boolean validate() {
        return true;
    }
    public void vote() {
        // 添加节点并随机分配代币数量
        Random random = new Random();
        int totalTokens = 10000; // 总代币数量
        for (int i = 0; i < 100; i++) {
            int tokenAmount = 1 + random.nextInt(totalTokens / 10); // 保证每个节点至少拥有1个代币
            totalTokens -= tokenAmount;
            String nodeAddress = HashUtils.sha256(UUID.randomUUID().toString());
            Node node = new Node(nodeAddress, tokenAmount);
            addNode(node);
            // 在添加节点的同时，创建对应的投票并添加到投票列表
            Vote vote = new Vote(nodeAddress, 0);
            addVote(vote);
            System.out.println("节点已添加，节点为："+ (i + 1) + ". " + node.getAddress() + "，代币数量为：" + node.getTokenAmount());
        }
        // 根据分配的代币给予节点票数
        for (Node node : nodeList) {
            int numVotes = node.getTokenAmount(); // 获取节点的代币数量
            node.addVote(numVotes); // 给节点增加票数
        }
        // 进行随机投票模拟
        Random random1 = new Random(System.currentTimeMillis());
        for (Node node : nodeList) {
            int numVotes = node.getVoteCount(); // 获取节点的票数
            for (int i = 0; i < numVotes; i++) {
                int candidateIndex = random1.nextInt(voteList.size()); // 随机选择候选人索引
                Vote vote = voteList.get(candidateIndex); // 获取对应的候选人投票
                node.vote(vote); // 节点进行投票
            }
        }
        for (Node node : nodeList) {
            for (Vote vote : voteList) {
                if (vote.getCandidate().equals(node.getAddress())) {
                    node.setVoteCount(vote.getVoteCount());
                    System.out.println(node.getVoteCount());
                }
            }
        }
        // 按票数排序节点
        List<Node> sortedNodes = sortNodesByVoteCount();

        // 输出票数最高的30个节点
        System.out.println("票数最高的30个节点：");
        for (int i = 0; i < 30 && i < sortedNodes.size(); i++) {
            Node node = sortedNodes.get(i);
            System.out.println((i + 1) + ". " + node.getAddress() + " - 票数：" + node.getVoteCount());
        }
        // 创建一个新的区块并添加到区块链
        Block newBlock1 = new Block(1, System.currentTimeMillis(), new ArrayList<>(), voteList, blockList.get(blockList.size()-1).hash, 0);
        System.out.println("等待添加区块1：");
        addBlock(newBlock1);
        System.out.println("区块 1已添加，区块哈希为：" + newBlock1.hash);
        Blockchain.Block newBlock2 = new Blockchain.Block(2, System.currentTimeMillis(), new ArrayList<>(), voteList, blockList.get(blockList.size()-1).hash, 0);
        System.out.println("等待添加区块2：");
        addBlock(newBlock2);
        System.out.println("区块 2已添加，区块哈希为：" + newBlock2.hash);
        // 验证区块链的合法性 
        System.out.println("区块链的合法性为：" + (validate() ? "validate" : "invalidate"));
    }
}
