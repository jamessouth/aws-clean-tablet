const path = require('path');
const webpack = require('webpack');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const ScriptExtHTMLWebpackPlugin = require('script-ext-html-webpack-plugin');
const TerserWebpackPlugin = require('terser-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;


module.exports = env => {
  const envObj = Object.keys(env)
    .reduce((acc, val) => {
      acc[`process.env.${val}`] = JSON.stringify(env[val]);
      return acc;
    }, {});

  const plugins = [
    new webpack.DefinePlugin(envObj),
    new webpack.DefinePlugin({
      'process.env.CT_WS': JSON.stringify(process.env.CT_WS)
    }),
    new CleanWebpackPlugin(),
    new MiniCssExtractPlugin({
      filename: '[name].[contenthash].css',
      chunkFilename: '[name].[contenthash].css',
    }),
    new HTMLWebpackPlugin({
      template: './src/front/index.html',
      title: 'Clean Tablet',
      favicon: './src/assets/icons/favicon-16x16.png'
    }),
    new ScriptExtHTMLWebpackPlugin({
      defaultAttribute: 'async',
    }),
    new webpack.HashedModuleIdsPlugin(),
  ];
  
  if (env.ENV == 'prod') {
    plugins.push(new BundleAnalyzerPlugin({
      analyzerMode: 'disabled',
      generateStatsFile: true,
    }));
  }

  return {
    mode: env.ENV == 'prod' ? 'production' : 'development',
    devtool: env.ENV == 'prod' ? false : 'source-map',
    entry: {
      main: './src/front/index.js',
    },
    output: {
      filename: env.ENV == 'prod' ? '[name].[contenthash].js' : '[name].[hash].js',
      chunkFilename: env.ENV == 'prod' ? '[name].[contenthash].js' : '[name].[hash].js',
      path: path.resolve(__dirname, 'dist'),
    },
    module: {
      rules: [
        {
          test: /\.m?js$/,
          include: path.resolve(__dirname, 'src/front/'),
          exclude: /(node_modules|\.test\.js$)/,
          use: {
            loader: 'babel-loader',
            options: {
              // presets: [
              //   [
              //     '@babel/preset-env',
              //     {
              //       'useBuiltIns': 'usage',
              //       'corejs': '3',
              //     },
              //   ],
              //   [
              //     '@babel/preset-react',
              //     {
              //       'useBuiltIns': true,
              //       'development': env.ENV == "dev",
              //     },
              //   ],
              // ],
              // plugins: [
              //   '@babel/plugin-transform-modules-commonjs'
              // ],
              cacheDirectory: true,
            },
          },
        },
        {
          test: /\.css$/i,
          include: path.resolve(__dirname, 'src'),
          exclude: /node_modules/,
          use: [
            // env.ENV == "dev" ?
            //   {
            //     loader: "style-loader",
            //     options: {
            //       esModule: false,
            //     },
            //   } :
              MiniCssExtractPlugin.loader,
            {
              loader: "css-loader",
              options: {
                sourceMap: env.ENV == "dev",
                importLoaders: 1,
              },
            },
            {
                loader: 'postcss-loader',
                options: {
                    sourceMap: env.ENV == "dev",
                },
            },
          ],
        },
        {
          test: /\.(png|svg|jpg|jpeg|gif)$/,
          use: [
            {
              loader: 'file-loader',
              options: {
                outputPath: 'images/',
                publicPath: 'images/',
              },
            },
          ],
        },
      ],
    },
    optimization: {
      minimizer: [
        new TerserWebpackPlugin({
          parallel: true,
          sourceMap: env.ENV == "dev",
        }),
      ],
      runtimeChunk: 'single',
      splitChunks: {
        chunks: 'all',
      },
    },
    plugins,
    devServer: {
      port: 4200,
      contentBase: path.join(__dirname, 'dist'),
      index: 'index.html',
      historyApiFallback: true,
      // hot: true,
    },
  }
};