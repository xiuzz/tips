public class Vote {

    private String candidate; // 候选人
    private int voteCount; // 获得票数

    public Vote(String candidate, int voteCount) {
        this.candidate = candidate;
        this.voteCount = voteCount;
    }

    public String getCandidate() {
        return this.candidate;
    }

    public int getVoteCount() {
        return this.voteCount;
    }

    public void incCount() {
        this.voteCount++;
    }
}
