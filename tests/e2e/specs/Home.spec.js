describe('Test home page', () => {
  it('Visits the app root url and sees the latest tracks', () => {
    cy.visit('/');

    cy.contains('h1', 'Latest songs');
    cy.contains('a', 'Details');
    cy.contains('a', 'to Spotify');
  });
});
