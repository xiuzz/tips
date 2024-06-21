public class BlockChain {
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
}
