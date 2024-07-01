public class Node {
    private String address;
    private int availableVotes; // 节点可用于投票的票数
    private int voteCount; // 获得票数
    private int tokenAmount; // 代币数量
    public Node(String address, int tokenAmount) {
        this.address = address;
        this.tokenAmount = tokenAmount;
    }
    public String getAddress() {
        return this.address;
    }
    public int getTokenAmount() {
        return this.tokenAmount;
    }

    public int getVoteCount() {
        return this.voteCount;
    }
    public void addVote(int votes) {
        this.availableVotes += votes;
    }

    public void vote(Vote vote) {
        if (this.availableVotes > 0) {
            vote.incCount();
            this.availableVotes--;
        }
    }

    public void setVoteCount(int count) {
        this.voteCount = count;
    }
}
