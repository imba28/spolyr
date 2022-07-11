const visitSearchPage = () => {
  cy.visit('/');
  cy.get('[aria-label="Search"]').click();
};

describe('Search page', () => {
  beforeEach(visitSearchPage);

  it('Searching for tracks should return results', () => {
    cy.contains(/\d+ tracks found/);
    cy.contains('Title');
    cy.contains('Detail');
  });

  it('Clicking the details button shows information about the track', async () => {
    cy.contains('Detail').first().click();

    cy.url().should('include', '/tracks/');
    cy.contains('Back');
    cy.contains('to Spotify');
  });
});
