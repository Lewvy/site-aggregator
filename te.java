class Demo {
    private static int count = 0;
    private int id;

    Demo() {
        count++;
        id = count;
    }
    static int getCount() {
        return count;
    }
    int getId() {
        return id;
    }
}

class Driver{
    public static void main(String args[]) {
       Demo d1 = new Demo();
       Demo d2 = new Demo();
       System.out.println(d1.getId());
       System.out.println(d2.getId());
       System.out.println(Demo.getCount());
    }
}
