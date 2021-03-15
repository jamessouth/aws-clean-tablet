class Bod {
    static pipe(...fns) {
      return function inner(start) {
        return fns.reduce((val, fn) => fn(val), start);
      };
    }

    static getDistAlongSide(max) {
        return Math.floor(Math.random() * (max + 1));
    }

    static getPoint(side, width, height) {
        // console.log(side, width, height);
        switch (side) {
          case 0:
            return [
              0,
              Bod.getDistAlongSide(height),
            ];
          case 1:
            return [
              Bod.getDistAlongSide(width),
              height,
            ];
          case 2:
            return [
              width,
              Bod.getDistAlongSide(height),
            ];
          case 3:
            return [
              Bod.getDistAlongSide(width),
              0,
            ];
          default: return undefined;
        }
      }
  
    static getNum(range, mod) {
      return Math.floor(Math.random() * range) - mod;
    }
  
    static getHypoLength(ang, height) {
      return height / Math.cos(ang); // 10 is the height of the border area
    }

    static getRandomPoint(width, height) {
        return [
          Math.floor(Math.random() * (width + 1)),
          Math.floor(Math.random() * (height + 1)),
        ];
      }
  
    static getCoord(hypo) {
      const dir = hypo < 0 ? -1 : 1;
      const opSide = Math.sqrt((hypo * hypo) - 100);
      return opSide * dir;
    }
  
    paint(ctx, {width, height}) { // eslint-disable-line


        // const ctr = [
        //     w / 2,
        //     h / 2,
        // ];

      for (let i = 0; i < 40; i += 1) {
        // const dir = Bod.getNum(Math.PI, 0);
  
        // const opLen = Bod.pipe(
        //   Bod.getHypoLength(h),
        //   Bod.getCoord,
        //   Math.round
        // )(dir);
  

        const [startx, starty] = Bod.getPoint(i%4, width, height);
        const [endx, endy] = Bod.getPoint((i+1)%4, width, height);
        console.log(startx, starty, endx, endy);
        const grad = ctx.createLinearGradient(startx, starty, endx, endy);
        grad.addColorStop(0, 'rgba(255, 76, 76, .07)');
        grad.addColorStop(0.12, 'transparent');
        grad.addColorStop(0.41, 'rgba(119, 150, 109, .08)');
        grad.addColorStop(0.62, 'transparent');
        grad.addColorStop(0.8, 'rgba(154, 148, 183, .09)');
        grad.addColorStop(1, 'rgba(189, 198, 103, .1)');
        ctx.fillStyle = grad;
        ctx.fillRect(0, 0, width, height);

        // ctx.beginPath();
        // ctx.moveTo(...getPoint(i % 4, w, h));
        // ctx.lineTo(stPt + opLen, 435);
        // ctx.lineWidth = Bod.getNum(10, -2);
        // ctx.strokeStyle = `hsl(${Bod.getNum(41, -317)}deg, ${Bod.getNum(30, -70)}%, ${Bod.getNum(30, -30)}%)`;
        // ctx.stroke();
      }
    }
  }
  registerPaint('bod', Bod); // eslint-disable-line