module.exports = {
  'env': {
    'browser': true,
    'es2021': true,
    'node': true,
    'jest/globals': true,
  },
  'plugins': ['jest'],
  'extends': [
    'eslint:recommended',
    'google',
    'plugin:vue/recommended',
  ],
  'parserOptions': {
    'ecmaVersion': 12,
    'sourceType': 'module',
  },
  'rules': {
    'max-len': ['error', {'code': 120}],
    'no-debugger': [process.env.NODE_ENV === 'production' ? 'error' : 'off'],
    'jest/no-disabled-tests': 'warn',
    'jest/no-focused-tests': 'error',
    'jest/no-identical-title': 'error',
    'jest/prefer-to-have-length': 'warn',
    'jest/valid-expect': 'error',
  },
  'ignorePatterns': ['assets/openapi/**/*.js', 'jest.config.js', 'vue.config.js'],
};
