// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
    config.set({
        basePath: '',
        frameworks: ['jasmine', '@angular/cli'],
        plugins: [
            require('karma-jasmine'),
            require('karma-chrome-launcher'),
            require('karma-coverage-istanbul-reporter'),
            require('karma-mocha-reporter'),
            require('karma-remap-istanbul'),
            require('karma-coverage'),
            require('@angular/cli/plugins/karma')
        ],
        files: [
            {pattern: './src/test.ts', watched: false}
        ],
        preprocessors: {
            "./src/test.ts": ["@angular/cli", "coverage"]
        },
        mime: {
            'text/x-typescript': ['ts', 'tsx']
        },
        remapIstanbulReporter: {
            reports: {
                html: 'coverage',
                lcovonly: './coverage/coverage.lcov'
            }
        },
        coverageIstanbulReporter: {
            reports: ['html', 'lcovonly', 'text-summary'],
            fixWebpackSourcePaths: true
        },
        angularCli: {
            config: './angular-cli.json',
            environment: 'dev'
        },
        reporters: config.angularCli && config.angularCli.codeCoverage
            ? ['mocha', 'karma-remap-istanbul']
            : ['mocha'],
        port: 9876,
        colors: true,
        logLevel: config.LOG_INFO,
        autoWatch: true,
        browsers:['Chrome_without_sandbox'],
        customLaunchers:{
            Chrome_without_sandbox:{
                base:'ChromeHeadless',
                flags:['-no-sandbox']
            }
        },
        singleRun: false
    });
};
