import {fireEvent, render} from '@testing-library/vue';

import SearchForm from '@/components/SearchForm';

describe('SearchForm', () => {
  test('Emits submit event when form is submitted', async () => {
    const {getByRole, emitted} = render(SearchForm);

    await fireEvent.click(getByRole('button'));

    expect(emitted().search).toHaveLength(1);
  });

  test('Emitted event contains the entered input value', async () => {
    const {getByRole, emitted} = render(SearchForm);

    const query = 'A simple query';
    await fireEvent.update(getByRole('textbox'), query);
    await fireEvent.click(getByRole('button'));

    expect(emitted().search[0][0]).toEqual(query);
  });
});
