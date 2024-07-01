public class Main {
    public static void main(String[] args) {
        Blockchain blockchain = new Blockchain();
        blockchain.createGenesisBlock();
        blockchain.vote();
    }
}
