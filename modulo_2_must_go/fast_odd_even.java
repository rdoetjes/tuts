import java.time.*;
import java.time.temporal.ChronoUnit;
import java.util.*;

class fast_odd_even {

  public static void andMethod(final int ITS) {
    int temp = 0;

    for (int i = 0; i < ITS; i++) {
      if ((i & 1) == 0) {
        temp++;
      }
    }
    System.out.printf("%d\n", temp);
  }

  public static void moduloMethod(final int ITS) {
    int temp = 0;

    for (int i = 0; i < ITS; i++) {
      if (i % 2 == 0) {
        temp++;
      }
    }
    System.out.printf("%d\n", temp);
  }

  private static double Average(List<Long> l) {
    Long sum = 0L;
    if (!l.isEmpty()) {
      for (Long e : l) {
        sum += e;
      }
      return sum.doubleValue() / l.size();
    }
    return sum;
  }

  public static void main(String[] args) {
    List<Long> aC = new ArrayList<>();
    List<Long> bC = new ArrayList<>();
    final int ITS = 10000000;

    Long start;
    Long stop;

    for (int i = 0; i < 100; i++) {
      start = ChronoUnit.MICROS.between(Instant.EPOCH, Instant.now());
      andMethod(ITS);
      stop = ChronoUnit.MICROS.between(Instant.EPOCH, Instant.now());
      aC.add(stop - start);

      start = ChronoUnit.MICROS.between(Instant.EPOCH, Instant.now());
      andMethod(ITS);
      stop = ChronoUnit.MICROS.between(Instant.EPOCH, Instant.now());
      bC.add(stop - start);
    }

    double aAvg = Average(aC);
    double bAvg = Average(bC);

    System.out.printf("and    avg: %.2f uSec\n", aAvg);
    System.out.printf("modulo avg: %.2f uSec\n", bAvg);
  }
}
