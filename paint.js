class Bod {
  static get inputProperties() { return ['--col1', '--col2', '--col3']; }
  static getDistAlongSide(max) {
    return Math.floor(Math.random() * (max + 1));
  }
  static getPoint(side, width, height) {
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
  paint(ctx, { width, height }, props) { // eslint-disable-line
    for (let i = 0; i < 14; i += 1) {
      const [startx, starty] = Bod.getPoint(i % 4, width, height);
      const [endx, endy] = Bod.getPoint((i + 1) % 4, width, height);
      const grad = ctx.createLinearGradient(startx, starty, endx, endy);
      grad.addColorStop(0, props.get('--col1'));
      grad.addColorStop(0.15, 'transparent');
      grad.addColorStop(0.58, props.get('--col1'));
      grad.addColorStop(0.71, props.get('--col2'));
      grad.addColorStop(0.79, 'transparent');
      grad.addColorStop(0.90, 'transparent');
      grad.addColorStop(1, props.get('--col3'));
      ctx.globalCompositeOperation = 'color-burn';
      ctx.fillStyle = grad;
      ctx.fillRect(0, 0, width, height);
    }
  }
}
registerPaint('bod', Bod); // eslint-disable-line