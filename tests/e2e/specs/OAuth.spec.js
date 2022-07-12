describe('Oauth smoke test', () => {
  it('Redirects to the Spotify oAuth2 sign in form', () => {
    cy.visit('/');

    cy.contains('Sign in').click();
    cy.url().should('include', 'https://accounts.spotify.com');
  });
});
